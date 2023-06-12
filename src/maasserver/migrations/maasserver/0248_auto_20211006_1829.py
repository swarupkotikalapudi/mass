# Generated by Django 2.2.12 on 2021-10-06 18:29

from django.db import migrations, models
import django.db.models.deletion

import maasserver.models.node


class Migration(migrations.Migration):
    dependencies = [
        ("maasserver", "0247_auto_20210915_1545"),
    ]

    operations = [
        migrations.AddField(
            model_name="vmcluster",
            name="pool",
            field=models.ForeignKey(
                blank=True,
                default=None,
                null=True,
                on_delete=django.db.models.deletion.PROTECT,
                to="maasserver.ResourcePool",
            ),
        ),
        migrations.AddField(
            model_name="vmcluster",
            name="zone",
            field=models.ForeignKey(
                default=maasserver.models.node.get_default_zone,
                on_delete=django.db.models.deletion.SET_DEFAULT,
                to="maasserver.Zone",
                verbose_name="Physical zone",
            ),
        ),
    ]