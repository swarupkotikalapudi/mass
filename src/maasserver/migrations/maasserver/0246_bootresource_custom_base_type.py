# Generated by Django 2.2.12 on 2021-08-24 14:35

from django.db import migrations, models


class Migration(migrations.Migration):
    dependencies = [
        ("maasserver", "0245_bmc_power_parameters_index_hash"),
    ]

    operations = [
        migrations.AddField(
            model_name="bootresource",
            name="base_image",
            field=models.CharField(blank=True, max_length=255),
        ),
    ]